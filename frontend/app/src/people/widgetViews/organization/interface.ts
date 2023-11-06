import { BountyRoles, BudgetHistory, Organization, PaymentHistory, Person } from 'store/main';

export interface ModalProps {
  isOpen: boolean;
  close: () => void;
  uuid?: string;
  user?: Person;
  addToast?: (text: string, color: 'danger' | 'success') => void;
}

export interface EditOrgModalProps extends ModalProps {
  org?: Organization;
  onSubmit: (body: any) => void;
  onDelete: () => void;
}

export interface UserRolesModalProps extends ModalProps {
  bountyRolesData: BountyRoles[];
  userRoles: any[];
  roleChange: (e: any, s: any) => void;
  submitRoles: () => void;
}

export interface PaymentHistoryModalProps extends ModalProps {
  url: string;
  paymentsHistory: PaymentHistory[];
}

export interface BudgetHistoryModalProps extends ModalProps {
  budgetsHistory: BudgetHistory[];
}

export interface AddUserModalProps extends ModalProps {
  loading: boolean;
  onSubmit: (body: any) => void;
  disableFormButtons: boolean;
  setDisableFormButtons: React.Dispatch<React.SetStateAction<boolean>>;
}

export interface AddBudgetModalProps extends ModalProps {
  invoiceStatus: boolean;
  startPolling: (inv: string) => void;
}

export interface WithdrawModalProps extends ModalProps {
  getOrganizationBudget: () => Promise<void>;
}
